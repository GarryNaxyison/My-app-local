import runpy
import sqlite3
import unittest
from datetime import datetime, timedelta, timezone
from pathlib import Path

report = runpy.run_path(str(Path(__file__).with_name("premium-transfer-report.py")))["report"]


class PremiumTransferReportTest(unittest.TestCase):
    def test_uses_mature_cohorts_and_bounded_outcomes(self):
        start = datetime(2026, 11, 1, tzinfo=timezone.utc)
        with sqlite3.connect(":memory:") as db:
            db.execute("CREATE TABLE learning_events (telegram_id, event_id, stage, format, created_at)")
            for user, event, stage, format_, day in (
                (1, "premium-transfer-v1", "premium_transfer_control", "", 0),
                (2, "premium-transfer-v1", "premium_transfer_offer", "", 1),
                (3, "premium-transfer-v1", "premium_transfer_offer", "", 0),
                (1, "purchase", "premium_transfer_purchased", "", 3),
                (1, "return", "today_opened", "", 7),
                (1, "retry", "retry_submitted", "voice", 2),
                (3, "late-purchase", "premium_transfer_purchased", "", 7),
                (3, "too-early", "today_opened", "", 6),
                (3, "too-late", "lesson_started", "", 8),
            ):
                db.execute("INSERT INTO learning_events VALUES (?, ?, ?, ?, ?)",
                           (user, event, stage, format_, (start + timedelta(days=day)).isoformat()))
            result = report(db, start + timedelta(days=8))
        self.assertEqual(result["control"], {
            "learners": 1, "purchasers": 1, "day_seven_returns": 1,
            "recorded_voice_retries": 1, "purchase_rate": 1, "day_seven_return_rate": 1,
        })
        self.assertEqual(result["offer"]["learners"], 1)
        self.assertEqual(result["offer"]["purchase_rate"], 0)
        self.assertEqual(result["offer"]["day_seven_return_rate"], 0)
