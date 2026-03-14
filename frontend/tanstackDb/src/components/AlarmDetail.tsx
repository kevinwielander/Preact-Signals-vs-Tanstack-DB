import { useAlarm } from "../hooks";
import EditText from "./EditText";
import EditSelect from "./EditSelect";
import EditList from "./EditList";

interface Props {
  alarmId: string | null;
}

export default function AlarmDetail({ alarmId }: Props) {
  const { data: alarm, isLoading } = useAlarm(alarmId);

  if (!alarmId) {
    return <div className="detail-placeholder">Select an alarm to view details</div>;
  }

  if (isLoading || !alarm) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <>
      <h2>Alarm Details</h2>
      <div className="detail-fields">
        <div className="field-row">
          <label>Title</label>
          <EditText alarmId={alarm.id} field="title" value={alarm.title} />
        </div>
        <div className="field-row">
          <label>Description</label>
          <EditText
            alarmId={alarm.id}
            field="description"
            value={alarm.description}
            multiline
          />
        </div>
        <div className="field-row">
          <label>Severity</label>
          <EditSelect
            alarmId={alarm.id}
            field="severity"
            value={alarm.severity}
            options={["low", "medium", "high", "critical"]}
          />
        </div>
        <div className="field-row">
          <label>Status</label>
          <EditSelect
            alarmId={alarm.id}
            field="status"
            value={alarm.status}
            options={["active", "acknowledged", "resolved"]}
          />
        </div>
        <div className="field-row">
          <label>Assigned Resources</label>
          <EditList items={alarm.assignedResources} />
        </div>
      </div>
      <div className="detail-meta">
        <div>
          Event # <span>{alarm.eventNumber}</span>
        </div>
        <div>
          Hash <span>{alarm.eventHash}</span>
        </div>
        <div>
          Created Offset <span>{alarm.createdOnOffset}</span>
        </div>
        <div>
          Archived Offset <span>{alarm.archivedOnOffset ?? "—"}</span>
        </div>
        <div>
          Created <span>{new Date(alarm.createdAt).toLocaleString()}</span>
        </div>
        <div>
          Updated <span>{new Date(alarm.updatedAt).toLocaleString()}</span>
        </div>
      </div>
    </>
  );
}
