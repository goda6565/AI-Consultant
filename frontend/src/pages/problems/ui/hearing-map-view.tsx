import type { HearingMap } from "@/shared/api";
import { Button, Loading, MermaidMindMap } from "@/shared/ui";

type HearingMapViewProps = {
  hearingMap?: HearingMap;
  isLoading: boolean;
  onCopyHearingMap: () => void;
};

export function HearingMapView({
  hearingMap,
  isLoading,
  onCopyHearingMap,
}: HearingMapViewProps) {
  return (
    <div className="w-full h-full p-8 rounded-xl border mb-8">
      <div className="flex justify-end items-center">
        <Button variant="outline" onClick={onCopyHearingMap}>
          ヒアリングマップをコピー
        </Button>
      </div>
      {isLoading && <Loading />}
      {!isLoading && hearingMap && (
        <div className="mx-auto">
          <MermaidMindMap chart={hearingMap?.content ?? ""} />
        </div>
      )}
    </div>
  );
}
