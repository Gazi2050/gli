import { ImageResponse } from "next/og";

export const config = {
  runtime: "edge",
};

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url);

    const title = searchParams.get("title") || "gli";
    const description = searchParams.get("description");

    return new ImageResponse(
      <div
        style={{
          height: "100%",
          width: "100%",
          display: "flex",
          flexDirection: "column",
          alignItems: "flex-start",
          justifyContent: "space-between",
          backgroundColor: "#0a0a0a",
          padding: "60px",
          color: "#22d3ee",
          fontFamily: "JetBrains Mono, monospace",
        }}
      >
        <h1 tw="text-7xl font-bold tracking-tight text-[#22d3ee] m-0 leading-tight">{title}</h1>
        {description ? <p tw="text-2xl text-[#a3a3a3] mt-6 max-w-4xl">{description}</p> : null}

        <div tw="flex flex-col w-full">
          <div tw="w-full h-px bg-[#262626] mb-[48px]" />
          <div tw="flex justify-between items-center w-full">
            <p> </p>
            <div tw="flex items-center">
              <span tw="text-3xl font-semibold text-[#a3a3a3]">gli - Git Wrapper</span>
            </div>
          </div>
        </div>
      </div>,
      {
        width: 1200,
        height: 630,
      }
    );
  } catch (e: any) {
    console.log(`${e.message}`);
    return new Response(`Failed to generate the image`, {
      status: 500,
    });
  }
}