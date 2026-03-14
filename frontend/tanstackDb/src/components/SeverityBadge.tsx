interface Props {
  severity: string;
}

export default function SeverityBadge({ severity }: Props) {
  return <span className={`severity-badge ${severity}`}>{severity}</span>;
}
