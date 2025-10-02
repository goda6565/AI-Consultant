import { useEffect, useState } from "react";
import type { ListHearingMessagesSuccessResponse } from "@/shared/api";
import { useListHearingMessages } from "@/shared/api";
import type { Message } from "../model/zod";

type UseChatMessageProps = {
  hearingId: string;
  enabled: boolean;
};

export const useChatMessage = ({ hearingId, enabled }: UseChatMessageProps) => {
  const [localMessages, setLocalMessages] = useState<Message[]>([]);

  const {
    data: hearingMessages,
    isLoading: isHearingMessagesLoading,
    error: isHearingMessagesError,
  } = useListHearingMessages(hearingId, {
    swr: {
      enabled: Boolean(hearingId) && enabled,
      refreshInterval: (data) =>
        enabled &&
        Boolean(hearingId) &&
        ((data as ListHearingMessagesSuccessResponse)?.hearingMessages
          ?.length ?? 0) === 0
          ? 1000
          : 0,
      revalidateOnFocus: false,
    },
  });

  useEffect(() => {
    if (hearingMessages) {
      setLocalMessages(hearingMessages.hearingMessages);
    }
  }, [hearingMessages]);

  const isLoading = isHearingMessagesLoading;
  const error = isHearingMessagesError;

  return {
    localMessages,
    setLocalMessages,
    isLoading,
    error,
  };
};
