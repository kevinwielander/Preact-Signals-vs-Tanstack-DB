interface Props {
  items: string[];
}

export default function EditList({ items }: Props) {
  if (items.length === 0) {
    return <span style={{ color: "#aaa" }}>None</span>;
  }

  return (
    <div className="tag-list">
      {items.map((item) => (
        <span key={item} className="tag">
          {item.slice(0, 8)}...
        </span>
      ))}
    </div>
  );
}
