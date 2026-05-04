import { source } from "@/loaders/source";
import defaultMdxComponents from "fumadocs-ui/mdx";
import { DocsBody, DocsPage, DocsTitle } from "fumadocs-ui/page";
import { notFound } from "next/navigation";

export const revalidate = 86400;
export const dynamic = "force-static";

export default async function Page(props: {
  params: Promise<{ slug?: string[] }>;
}) {
  const params = await props.params;
  const page = source.getPage(params.slug);
  if (!page) notFound();

  const title = page.data.title;
  const MDXContent = page.data.body;
  const toc = page.data.toc.filter((item) => item.depth <= 3);

  return (
    <DocsPage toc={toc} tableOfContent={{ style: "clerk", single: false }} full={false}>
      {title && title !== "Intro" && (
        <div className="mb-6">
          <DocsTitle>{title}</DocsTitle>
        </div>
      )}
      <DocsBody>
        <MDXContent components={defaultMdxComponents} />
      </DocsBody>
    </DocsPage>
  );
}

export async function generateStaticParams() {
  return source.generateParams();
}

export async function generateMetadata(props: {
  params: Promise<{ slug?: string[] }>;
}) {
  const params = await props.params;
  const page = source.getPage(params.slug);
  if (!page) notFound();

  const rootTitle = page.data.title ?? "Home";
  const title = rootTitle + " | gli";
  const description = page.data.description;
  return {
    title,
    description,
  };
}