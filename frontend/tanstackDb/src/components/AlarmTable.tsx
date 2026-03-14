import { useState } from "react";
import { useAlarmsPaginated } from "../hooks";
import SeverityBadge from "./SeverityBadge";

interface Props {
  selectedId: string | null;
  onSelectAlarm: (id: string) => void;
}

const PAGE_SIZE = 20;

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
            <tr
              key={alarm.id}
              className={alarm.id === selectedId ? "selected" : ""}
              onClick={() => onSelectAlarm(alarm.id)}
            >
              <td>{alarm.title}</td>
              <td>
                <SeverityBadge severity={alarm.severity} />
              </td>
              <td>{alarm.status}</td>
              <td>{new Date(alarm.createdAt).toLocaleString()}</td>
            </tr>
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
