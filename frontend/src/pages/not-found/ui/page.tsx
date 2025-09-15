import { Button, Heading, RegularText } from "@/shared/ui";

export function NotFound() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center space-y-6 text-center">
      <div className="space-y-4">
        <div className="text-6xl font-bold">404</div>
        <Heading>Not Found</Heading>
        <RegularText>
          お探しのページは存在しないか、移動または削除された可能性があります。
        </RegularText>
      </div>

      <div className="flex space-x-4">
        <Button asChild>
          <a href="/">ホームに戻る</a>
        </Button>
      </div>
    </div>
  );
}
