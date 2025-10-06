import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { HearingMapView } from "@/pages/problems/ui/hearing-map-view";
import type { HearingMap } from "@/shared/api";

const meta: Meta<typeof HearingMapView> = {
  title: "Pages/Problems/HearingMapView",
  component: HearingMapView,
};

export default meta;

type Story = StoryObj<typeof HearingMapView>;

const sampleHearingMap: HearingMap = {
  id: "hm-1",
  hearingId: "hearing-1",
  problemId: "problem-1",
  content: `mindmap
  root((顧客体験の低下))
    待ち時間
      ピーク帯集中
      手続き複雑
    デジタルUX
      アプリ利用率低
      案内不足
    従業員負荷
      研修不足
      事務作業多い`,
};

export const Default: Story = {
  args: {
    hearingMap: sampleHearingMap,
    isLoading: false,
    onCopyHearingMap: () => {
      console.log("Hearing map copied");
    },
  },
};

export const Loading: Story = {
  args: {
    hearingMap: undefined,
    isLoading: true,
    onCopyHearingMap: () => {
      console.log("Hearing map copied");
    },
  },
};

export const NoHearingMap: Story = {
  args: {
    hearingMap: undefined,
    isLoading: false,
    onCopyHearingMap: () => {
      console.log("Hearing map copied");
    },
  },
};
