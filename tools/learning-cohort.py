"""Print aggregate seven-day outcomes from the SQLite learning database."""

import argparse
import json
import sqlite3
from datetime import datetime, timedelta, timezone
from pathlib import Path


def cohort(rows, cutoff):
    first = {}
    for learner, corrected, completed in rows:
        correction = datetime.fromisoformat(corrected)
        if learner not in first or correction < first[learner][0]:
            first[learner] = (
                correction,
                datetime.fromisoformat(completed) if completed else None,
            )
    eligible = [entry for entry in first.values() if entry[0] <= cutoff]
    successful = sum(
        end is not None and start <= end <= start + timedelta(days=7)
        for start, end in eligible
    )
    return {
        "eligible_learners": len(eligible),
        "successful_learners": successful,
        "completion_rate": successful / len(eligible) if eligible else None,
    }


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("database", type=Path)
    parser.add_argument("--language", default="en")
    args = parser.parse_args()
    cutoff = datetime.now(timezone.utc) - timedelta(days=7)
    with sqlite3.connect(args.database.resolve().as_uri() + "?mode=ro", uri=True) as db:
        rows = db.execute(
            """SELECT telegram_id, corrected_at, completed_at
               FROM mistake_transfers
               WHERE language = ? AND corrected_at <> ''
               ORDER BY telegram_id, review_id""",
            (args.language,),
        )
        print(json.dumps(cohort(rows, cutoff)))


if __name__ == "__main__":
    main()
