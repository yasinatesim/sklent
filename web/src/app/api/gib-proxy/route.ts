import { NextResponse } from "next/server";
import { isGibRequest } from "@/lib/gib-proxy";

export const GET = async (request: Request) => {
  const { searchParams } = new URL(request.url);
  const target = searchParams.get("target") ?? "";

  if (!isGibRequest(target)) {
    return NextResponse.json({ error: "host_not_allowed" }, { status: 400 });
  }

  const upstream = await fetch(target, { headers: { Accept: "application/json" } });
  const body = await upstream.text();
  return new NextResponse(body, {
    status: upstream.status,
    headers: { "Content-Type": upstream.headers.get("Content-Type") ?? "text/plain" },
  });
};
