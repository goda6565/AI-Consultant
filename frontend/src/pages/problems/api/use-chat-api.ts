import { useGetHearing, useGetProblem } from "@/shared/api";

export const useChatApi = (problemId: string) => {
  const { data: problem, isLoading, error, mutate } = useGetProblem(problemId);
  const {
    data: hearing,
    isLoading: isHearingLoading,
    error: isHearingError,
    mutate: mutateHearing,
  } = useGetHearing(problemId);

  const mutateChat = () => {
    mutate();
    mutateHearing();
  };

  const isChatLoading = isLoading || isHearingLoading;
  const isChatError = error || isHearingError;

  return {
    mutateChat,
    isChatLoading,
    isChatError,
    problem,
    hearing,
  };
};
