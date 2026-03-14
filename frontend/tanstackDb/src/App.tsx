import { useState } from "react";
import { useMe } from "./hooks";
import { useWebSocket } from "./useWebSocket";
import Sidebar from "./components/Sidebar";
import AlarmDetail from "./components/AlarmDetail";
import AlarmTable from "./components/AlarmTable";

export default function App() {
  const [selectedAlarmId, setSelectedAlarmId] = useState<string | null>(null);
  const { data: me } = useMe();
  useWebSocket();

  return (
    <div className="app-layout">
      <aside className="sidebar">
        <Sidebar
          meId={me?.id}
          selectedId={selectedAlarmId}
          onSelectAlarm={setSelectedAlarmId}
        />
      </aside>
      <main className="detail">
        <AlarmDetail alarmId={selectedAlarmId} />
      </main>
      <footer className="alarm-table">
        <AlarmTable
          selectedId={selectedAlarmId}
          onSelectAlarm={setSelectedAlarmId}
        />
      </footer>
    </div>
  );
}
