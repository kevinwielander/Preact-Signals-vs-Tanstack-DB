import { useMyAlarms } from "../hooks";
import SeverityBadge from "./SeverityBadge";

interface Props {
  meId: string | undefined;
  selectedId: string | null;
  onSelectAlarm: (id: string) => void;
}

export default function Sidebar({ meId, selectedId, onSelectAlarm }: Props) {
  const { data, isLoading } = useMyAlarms(meId);

  return (
    <>
      <h2>My Alarms</h2>
      {isLoading && <div className="loading">Loading...</div>}
      {data?.items.length === 0 && (
        <div className="loading">No alarms assigned to you</div>
      )}
      {data?.items.map((alarm) => (
        <div
          key={alarm.id}
          className={`sidebar-item ${alarm.id === selectedId ? "selected" : ""}`}
          onClick={() => onSelectAlarm(alarm.id)}
        >
          <div className="alarm-title">{alarm.title}</div>
          <div className="alarm-meta">
            <SeverityBadge severity={alarm.severity} />
            <span>{alarm.status}</span>
          </div>
        </div>
      ))}
    </>
  );
}
