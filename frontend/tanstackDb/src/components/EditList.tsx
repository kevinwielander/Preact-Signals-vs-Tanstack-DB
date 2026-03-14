import { useResourceMap } from "../hooks";

interface Props {
  items: string[];
}

export default function EditList({ items }: Props) {
  const resourceMap = useResourceMap();

  if (items.length === 0) {
    return <span style={{ color: "#aaa" }}>None</span>;
  }

  return (
    <div className="tag-list">
      {items.map((id) => (
        <span key={id} className="tag">
          {resourceMap.get(id) ?? id.slice(0, 8) + "..."}
        </span>
      ))}
    </div>
  );
}
