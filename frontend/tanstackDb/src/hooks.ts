import {
  useQuery,
  useMutation,
  useQueryClient,
  keepPreviousData,
} from "@tanstack/react-query";
import type { Alarm } from "./types";
import { fetchMe, fetchAlarms, fetchAlarm, fetchResources, patchAlarm } from "./api";

export function useMe() {
  return useQuery({
    queryKey: ["me"],
    queryFn: fetchMe,
    staleTime: Infinity,
  });
}

export function useResources() {
  return useQuery({
    queryKey: ["resources"],
    queryFn: fetchResources,
    staleTime: Infinity,
  });
}

export function useResourceMap() {
  const { data: resources } = useResources();
  const map = new Map<string, string>();
  if (resources) {
    for (const r of resources) {
      map.set(r.id, r.displayName);
    }
  }
  return map;
}

export function useMyAlarms(meId: string | undefined) {
  const queryClient = useQueryClient();
  return useQuery({
    queryKey: ["alarms", "mine", meId],
    queryFn: async () => {
      const result = await fetchAlarms(1, 1000, meId!);
      for (const alarm of result.items) {
        queryClient.setQueryData(["alarm", alarm.id], alarm);
      }
      return result;
    },
    enabled: !!meId,
  });
}

export function useAlarmsPaginated(page: number, pageSize: number) {
  const queryClient = useQueryClient();
  return useQuery({
    queryKey: ["alarms", "paginated", page, pageSize],
    queryFn: async () => {
      const result = await fetchAlarms(page, pageSize);
      for (const alarm of result.items) {
        queryClient.setQueryData(["alarm", alarm.id], alarm);
      }
      return result;
    },
    placeholderData: keepPreviousData,
  });
}

export function useAlarm(id: string | null) {
  return useQuery({
    queryKey: ["alarm", id],
    queryFn: () => fetchAlarm(id!),
    enabled: !!id,
  });
}

export function useAlarmField<K extends keyof Alarm>(alarmId: string, field: K) {
  return useQuery({
    queryKey: ["alarm", alarmId],
    queryFn: () => fetchAlarm(alarmId),
    select: (alarm) => alarm[field],
    enabled: !!alarmId,
  });
}

export function usePatchAlarm() {
  const queryClient = useQueryClient();
  const { data: me } = useMe();
  return useMutation({
    mutationFn: ({
      id,
      field,
      value,
    }: {
      id: string;
      field: string;
      value: unknown;
    }) => patchAlarm(id, field, value, me?.id),
    onSuccess: (_data, { id }) => {
      queryClient.invalidateQueries({ queryKey: ["alarm", id] });
      queryClient.invalidateQueries({ queryKey: ["alarms"] });
    },
  });
}
