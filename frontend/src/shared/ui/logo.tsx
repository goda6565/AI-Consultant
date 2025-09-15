import Image from "next/image";

export function Logo({ className }: { className?: string }) {
  return (
    <a href="/">
      <Image
        src="/logo.png"
        alt="Logo"
        width={30}
        height={30}
        className={className}
      />
    </a>
  );
}
