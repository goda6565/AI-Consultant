import Image from "next/image";
import Link from "next/link";

export function Logo({ className }: { className?: string }) {
  return (
    <Link href="/">
      <Image
        src="/logo.png"
        alt="Logo"
        width={30}
        height={30}
        className={className}
      />
    </Link>
  );
}
