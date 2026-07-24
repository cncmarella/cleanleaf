/**
 * Base URL of the Go API. Public because the contact form posts from the
 * browser; set it in `.env.local` for dev and in Vercel for deploys.
 */
export const apiUrl =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export type ContactPayload = {
  name: string;
  email: string;
  phone: string;
  subject: string;
  message: string;
  source: string; // which page/CTA led here, for enquiry attribution
  website: string; // honeypot; always submitted empty by real users
};

export type ContactResult =
  | { ok: true; message: string }
  | { ok: false; message: string; fields?: Record<string, string> };

/** Mirrors the JSON error shape the Go API returns for every failure. */
type ApiError = { error?: string; fields?: Record<string, string> };

export async function submitContact(
  payload: ContactPayload,
  signal?: AbortSignal,
): Promise<ContactResult> {
  let response: Response;
  try {
    response = await fetch(`${apiUrl}/api/contact`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
      signal,
    });
  } catch {
    // Network failure, CORS rejection, or the API being down all land here.
    return {
      ok: false,
      message:
        "We could not reach our servers. Please check your connection or call us instead.",
    };
  }

  // A proxy returning HTML on error would otherwise blow up in .json().
  const body: ApiError = await response.json().catch(() => ({}));

  if (response.ok) {
    return { ok: true, message: "Thanks for reaching out. We'll be in touch shortly." };
  }

  return {
    ok: false,
    message: body.error ?? "Something went wrong. Please try again.",
    fields: body.fields,
  };
}
