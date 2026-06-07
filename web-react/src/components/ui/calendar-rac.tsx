"use client";

import { cn } from "@/lib/utils";
import { getLocalTimeZone, today } from "@internationalized/date";
import type { ComponentProps } from "react";
import {
  Button,
  Calendar as CalendarRac,
  CalendarCell as CalendarCellRac,
  CalendarGrid as CalendarGridRac,
  CalendarGridBody as CalendarGridBodyRac,
  CalendarGridHeader as CalendarGridHeaderRac,
  CalendarHeaderCell as CalendarHeaderCellRac,
  Heading as HeadingRac,
  RangeCalendar as RangeCalendarRac,
  composeRenderProps,
} from "react-aria-components";
import { ChevronLeftIcon, ChevronRightIcon } from "@radix-ui/react-icons";

interface BaseCalendarProps {
  className?: string;
}

type CalendarProps = ComponentProps<typeof CalendarRac> & BaseCalendarProps & {
  completedDates?: string[];
  loginDates?: string[];
};
type RangeCalendarProps = ComponentProps<typeof RangeCalendarRac> & BaseCalendarProps;

const dateKey = (date: { year: number; month: number; day: number }) =>
  `${date.year}-${String(date.month).padStart(2, "0")}-${String(date.day).padStart(2, "0")}`;

const CalendarHeader = () => (
  <header className="calendar-rac-header">
    <Button slot="previous" className="calendar-rac-nav">
      <ChevronLeftIcon width={16} height={16} />
    </Button>
    <HeadingRac className="calendar-rac-heading" />
    <Button slot="next" className="calendar-rac-nav">
      <ChevronRightIcon width={16} height={16} />
    </Button>
  </header>
);

const CalendarGridComponent = ({
  isRange = false,
  completedDates = [],
  loginDates = [],
}: {
  isRange?: boolean;
  completedDates?: string[];
  loginDates?: string[];
}) => {
  const now = today(getLocalTimeZone());
  const completed = new Set(completedDates);
  const logged = new Set(loginDates);

  return (
    <CalendarGridRac className="calendar-rac-grid">
      <CalendarGridHeaderRac>
        {(day) => (
          <CalendarHeaderCellRac className="calendar-rac-weekday">
            {day}
          </CalendarHeaderCellRac>
        )}
      </CalendarGridHeaderRac>
      <CalendarGridBodyRac className="calendar-rac-body">
        {(date) => {
          const key = dateKey(date);
          return (
            <CalendarCellRac
              date={date}
              data-login={logged.has(key) ? "true" : undefined}
              data-complete={completed.has(key) ? "true" : undefined}
              className={cn(
                "calendar-rac-cell",
                isRange && "calendar-rac-cell--range",
                date.compare(now) === 0 && "calendar-rac-cell--today",
              )}
            />
          );
        }}
      </CalendarGridBodyRac>
    </CalendarGridRac>
  );
};

const Calendar = ({ className, completedDates, loginDates, ...props }: CalendarProps) => {
  return (
    <CalendarRac
      {...props}
      className={composeRenderProps(className, (value) =>
        cn("calendar-rac", value),
      )}
    >
      <CalendarHeader />
      <CalendarGridComponent completedDates={completedDates} loginDates={loginDates} />
    </CalendarRac>
  );
};

const RangeCalendar = ({ className, ...props }: RangeCalendarProps) => {
  return (
    <RangeCalendarRac
      {...props}
      className={composeRenderProps(className, (value) =>
        cn("calendar-rac", value),
      )}
    >
      <CalendarHeader />
      <CalendarGridComponent isRange />
    </RangeCalendarRac>
  );
};

export { Calendar, RangeCalendar };
