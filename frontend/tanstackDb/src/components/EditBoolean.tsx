import { usePatchAlarm } from "../hooks";

interface Props {
  alarmId: string;
  field: string;
  value: boolean;
}

export default function EditBoolean({ alarmId, field, value }: Props) {
  const patch = usePatchAlarm();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    patch.mutate({ id: alarmId, field, value: e.target.checked });
  };

  return (
    <div className="checkbox-row">
      <input type="checkbox" checked={value} onChange={handleChange} />
    </div>
  );
}
