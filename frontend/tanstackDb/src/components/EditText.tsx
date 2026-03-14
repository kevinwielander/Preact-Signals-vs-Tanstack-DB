import { useState, useEffect } from "react";
import { usePatchAlarm } from "../hooks";

interface Props {
  alarmId: string;
  field: string;
  value: string;
  multiline?: boolean;
}

export default function EditText({ alarmId, field, value, multiline }: Props) {
  const [localValue, setLocalValue] = useState(value);
  const patch = usePatchAlarm();

  useEffect(() => {
    setLocalValue(value);
  }, [value]);

  const handleBlur = () => {
    if (localValue !== value) {
      patch.mutate({ id: alarmId, field, value: localValue });
    }
  };

  if (multiline) {
    return (
      <textarea
        value={localValue}
        onChange={(e) => setLocalValue(e.target.value)}
        onBlur={handleBlur}
      />
    );
  }

  return (
    <input
      type="text"
      value={localValue}
      onChange={(e) => setLocalValue(e.target.value)}
      onBlur={handleBlur}
    />
  );
}
