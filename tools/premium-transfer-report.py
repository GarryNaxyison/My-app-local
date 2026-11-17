"""Report aggregate Premium experiment outcomes from SQLite (read-only)."""

import argparse
import json
import sqlite3
from datetime import datetime, timezone
from pathlib import Path


def report(db, as_of):
    db.row_factory = sqlite3.Row
    rows = db.execute(
        """
        WITH enrolled AS (
          SELECT telegram_id, stage AS variant, julianday(created_at) AS start
          FROM learning_events
          WHERE event_id = 'premium-transfer-v1'
            AND stage IN ('premium_transfer_control', 'premium_transfer_offer')
            AND julianday(created_at) + 8 <= julianday(?)
        ), outcomes AS (
          SELECT enrolled.telegram_id, variant,
            MAX(CASE WHEN stage = 'premium_transfer_purchased'
              AND julianday(created_at) < start + 7 THEN 1 ELSE 0 END) AS purchased,
            MAX(CASE WHEN stage IN ('today_opened', 'lesson_started', 'answer_submitted',
              'answer_checked', 'retry_submitted', 'review_due', 'transfer_completed',
              'tutor_started', 'tutor_answer_submitted', 'tutor_answer_checked')
              AND julianday(created_at) >= start + 7 THEN 1 ELSE 0 END) AS returned,
            SUM(CASE WHEN stage = 'retry_submitted' AND format = 'voice'
              AND julianday(created_at) < start + 7 THEN 1 ELSE 0 END) AS voice_retries
          FROM enrolled LEFT JOIN learning_events ON learning_events.telegram_id = enrolled.telegram_id
            AND julianday(created_at) >= start AND julianday(created_at) < start + 8
          GROUP BY enrolled.telegram_id, variant
        )
        SELECT variant, COUNT(*) AS learners, SUM(purchased) AS purchasers,
          SUM(returned) AS day_seven_returns, SUM(voice_retries) AS recorded_voice_retries
        FROM outcomes GROUP BY variant
        """,
        (as_of.isoformat(),),
    )
    result = {variant: {
        "learners": 0, "purchasers": 0, "day_seven_returns": 0, "recorded_voice_retries": 0,
        "purchase_rate": None, "day_seven_return_rate": None,
    } for variant in ("control", "offer")}
    for row in rows:
        metrics = dict(row)
        variant = metrics.pop("variant").removeprefix("premium_transfer_")
        metrics["purchase_rate"] = row["purchasers"] / row["learners"]
        metrics["day_seven_return_rate"] = row["day_seven_returns"] / row["learners"]
        result[variant] = metrics
    return result


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("database", type=Path)
    args = parser.parse_args()
    with sqlite3.connect(args.database.resolve().as_uri() + "?mode=ro", uri=True) as db:
        print(json.dumps(report(db, datetime.now(timezone.utc)), indent=2))
