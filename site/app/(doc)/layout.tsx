import { baseOptions } from "@/app/layout.config";
import { source } from "@/loaders/source";
import { DocsLayout, type DocsLayoutProps } from "fumadocs-ui/layouts/docs";
import type { ReactNode } from "react";

export const dynamic = "force-static";

export const layoutProps: DocsLayoutProps = {
  ...baseOptions,
  tree: source.pageTree,
};

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <DocsLayout {...layoutProps} containerProps={{}}>
      {children}
    </DocsLayout>
  );
}