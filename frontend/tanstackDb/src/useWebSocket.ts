import { useEffect } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { Alarm } from "./types";

export function useWebSocket() {
  const queryClient = useQueryClient();

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws");

    ws.onmessage = (e) => {
      const event = JSON.parse(e.data);

      if (event.aggregateType === "alarm") {
        if (event.eventType === "AlarmFieldUpdated") {
          const data = event.data;
          queryClient.setQueryData<Alarm>(
            ["alarm", event.aggregateId],
            (old) => {
              if (!old) return old;
              return {
                ...old,
                [data.field]: data.newValue,
                eventHash: event.hash,
                eventNumber: event.version,
                updatedAt: event.timestamp,
              };
            }
          );
        } else if (event.eventType === "AlarmCreated") {
          queryClient.invalidateQueries({ queryKey: ["alarms"] });
        }
      }

      if (event.aggregateType === "resource") {
        queryClient.invalidateQueries({ queryKey: ["resources"] });
      }
    };

    return () => ws.close();
  }, [queryClient]);
}
