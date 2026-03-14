import { usePatchAlarm } from "../hooks";

interface Props {
  alarmId: string;
  field: string;
  value: string;
  options: string[];
}

export default function EditSelect({ alarmId, field, value, options }: Props) {
  const patch = usePatchAlarm();

  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    patch.mutate({ id: alarmId, field, value: e.target.value });
  };

  return (
    <select value={value} onChange={handleChange}>
      {options.map((opt) => (
        <option key={opt} value={opt}>
          {opt}
        </option>
      ))}
    </select>
  );
}
