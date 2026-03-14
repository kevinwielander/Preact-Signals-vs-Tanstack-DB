import { memo, useState } from "react";
import { useAlarmsPaginated, useAlarmField } from "../hooks";
import type { Alarm } from "../types";
import SeverityBadge from "./SeverityBadge";

interface Props {
  selectedId: string | null;
  onSelectAlarm: (id: string) => void;
}

const PAGE_SIZE = 20;

function TextCell({ alarmId, field }: { alarmId: string; field: keyof Alarm }) {
  const { data } = useAlarmField(alarmId, field);
  return <td>{String(data ?? "")}</td>;
}

function SeverityCell({ alarmId }: { alarmId: string }) {
  const { data } = useAlarmField(alarmId, "severity");
  return (
    <td>
      <SeverityBadge severity={data ?? "low"} />
    </td>
  );
}

function DateCell({ alarmId, field }: { alarmId: string; field: "createdAt" | "updatedAt" }) {
  const { data } = useAlarmField(alarmId, field);
  return <td>{data ? new Date(data).toLocaleString() : ""}</td>;
}

const AlarmRow = memo(function AlarmRow({
  alarmId,
  isSelected,
  onSelect,
}: {
  alarmId: string;
  isSelected: boolean;
  onSelect: (id: string) => void;
}) {
  return (
    <tr
      className={isSelected ? "selected" : ""}
      onClick={() => onSelect(alarmId)}
    >
      <TextCell alarmId={alarmId} field="title" />
      <SeverityCell alarmId={alarmId} />
      <TextCell alarmId={alarmId} field="status" />
      <DateCell alarmId={alarmId} field="createdAt" />
    </tr>
  );
});

export default function AlarmTable({ selectedId, onSelectAlarm }: Props) {
  const [page, setPage] = useState(1);
  const { data, isLoading } = useAlarmsPaginated(page, PAGE_SIZE);

  const totalPages = data ? Math.ceil(data.total / data.pageSize) : 0;

  return (
    <>
      <h2>All Alarms</h2>
      {isLoading && <div className="loading">Loading...</div>}
      <table>
        <thead>
          <tr>
            <th>Title</th>
            <th>Severity</th>
            <th>Status</th>
            <th>Created</th>
          </tr>
        </thead>
        <tbody>
          {data?.items.map((alarm) => (
            <AlarmRow
              key={alarm.id}
              alarmId={alarm.id}
              isSelected={alarm.id === selectedId}
              onSelect={onSelectAlarm}
            />
          ))}
        </tbody>
      </table>
      <div className="pagination">
        <button disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>
          Previous
        </button>
        <span>
          Page {page} of {totalPages} ({data?.total ?? 0} alarms)
        </span>
        <button
          disabled={page >= totalPages}
          onClick={() => setPage((p) => p + 1)}
        >
          Next
        </button>
      </div>
    </>
  );
}
