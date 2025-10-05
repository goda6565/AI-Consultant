import { useGetHearing, useGetJobConfig, useGetProblem } from "@/shared/api";

export const useChatApi = (problemId: string) => {
  const { data: problem, isLoading, error, mutate } = useGetProblem(problemId);
  const {
    data: hearing,
    isLoading: isHearingLoading,
    error: isHearingError,
    mutate: mutateHearing,
  } = useGetHearing(problemId);

  const {
    data: jobConfig,
    isLoading: isJobConfigLoading,
    error: isJobConfigError,
    mutate: mutateJobConfig,
  } = useGetJobConfig(problemId);

  const mutateChat = () => {
    mutate();
    mutateHearing();
    mutateJobConfig();
  };

  const isChatLoading = isLoading || isHearingLoading || isJobConfigLoading;
  const isChatError = error || isHearingError || isJobConfigError;

  return {
    mutateChat,
    isChatLoading,
    isChatError,
    problem,
    hearing,
    jobConfig,
  };
};
