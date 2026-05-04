import { source } from "@/loaders/source";

export const revalidate = false;

function stringifyTitle(title: any): string {
  if (typeof title === "string") return title;
  if (typeof title === "number" || typeof title === "bigint" || typeof title === "boolean") {
    return String(title);
  }

  if (title && typeof title === "object") {
    if (title.props?.children) {
      return stringifyTitle(title.props.children);
    }

    if (Array.isArray(title)) {
      return title.map(stringifyTitle).join("");
    }

    if (title.textContent) {
      return String(title.textContent);
    }

    if (title.type && (title.type === "span" || title.type === "code" || typeof title.type === "string")) {
      if (typeof title.props?.children === "string") {
        return title.props.children;
      }
    }
  }

  return "";
}

export async function GET() {
  const pages = source.getPages();

  let txt = `# gli

> gli is a modern Git wrapper that simplifies developer workflows.

`;

  for (const page of pages) {
    const title = stringifyTitle(page.data.title);
    if (title.startsWith("---")) {
      continue;
    }

    txt += `## ${title || "Untitled"}\n\n`;

    const pageUrl = `https://gli-tool.vercel.app${page.url}`;
    const description = page.data.description || "Documentation page";
    txt += `- [${title || "Untitled"}](${pageUrl}): ${description}\n`;

    if (page.data.toc && page.data.toc.length > 0) {
      const sections = page.data.toc.filter((item: any) => item.depth >= 2 && item.depth <= 4);

      if (sections.length > 0) {
        txt += "\n";

        for (const section of sections) {
          const sectionTitle = stringifyTitle(section.title);
          if (!sectionTitle) continue;

          const anchor = section.url.replace("#", "");
          const fullUrl = `https://gli-tool.vercel.app${page.url}?id=${anchor}`;

          txt += `- [${sectionTitle}](${fullUrl})\n`;
        }
      }
    }

    txt += "\n";
  }

  txt += `---

This documentation covers gli, a modern Git wrapper CLI. Use the URLs above to access specific pages and sections for detailed information about installation, commands, and configuration.
`;

  return new Response(txt, {
    headers: {
      "Content-Type": "text/plain; charset=utf-8",
    },
  });
}