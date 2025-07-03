"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";
import { FormEvent, useEffect, useState } from "react";
import { Header } from '@/components/Header'

export default function Validate() {
    const [email, setEmail] = useState<string | null>(null);
    const [code, setCode] = useState("");
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState("The six digit code has been sent to your email.");
    const [error, setError] = useState("");
    const router = useRouter();

    useEffect(() => {
        const storedEmail = localStorage.getItem("signin_email");
        setEmail(storedEmail);
    }, []);

    const handleResend = async () => {
        setLoading(true);
        setMessage("");

        try {
            const res = await fetch('/api/resend', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email }),
            });

            if (!res.ok) {
                throw new Error('Failed to resend code');
            }

            setMessage("Code sent!");
        } catch (err: unknown) {
            const errorMsg = err instanceof Error ? err.message : "An unexpected error occurred.";
            setError(errorMsg);
         } finally {
            setLoading(false);
            setMessage("The six digit code has been sent to your email again.")
        }
    };

    async function handleSubmit(e: FormEvent) {
        e.preventDefault();
        setLoading(true);
        setError("");

        const payload = {
            email,
            code,
        };

        try {
            const res = await fetch(
                `/api/verify`,
                {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    credentials: "include",
                    body: JSON.stringify(payload),
                },
            );
            if (!res.ok) throw new Error("Failed to sign in. Please try again.");
            router.push("https://app.mylocal.ing");
        } catch (err: unknown) {
            const message =
                err instanceof Error ? err.message : "An unexpected error occurred.";
            setError(message);
            setLoading(false);
        }
    }

    return (
        <>
            <Header />
            <div className="flex flex-col px-6 py-12 lg:px-8">
                <div className="sm:mx-auto sm:w-full sm:max-w-sm">
                    <Image
                        width={110}
                        height={51}
                        className="mx-auto h-30 w-auto"
                        src="/myLocal.svg"
                        alt="myLocal Logo"
                        priority
                    />
                </div>
                <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
                    <form onSubmit={handleSubmit} className="space-y-6">
                        {/* code */}
                        <div>
                            <label htmlFor="code" className="block text-sm/6 font-medium text-gray-900">
                                Code
                            </label>
                            <div className="mt-2">
                                <input
                                    id="code"
                                    type="text"
                                    required
                                    pattern="\d{6}"
                                    value={code}
                                    onChange={(e) => setCode(e.target.value)}
                                    className="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-yellow-600 sm:text-sm/6"
                                />
                            </div>
                        </div>

                        {/* submit */}
                        <div>
                            <button
                                type="submit"
                                disabled={loading}
                                className="flex w-full justify-center rounded-md bg-yellow-600 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-xs hover:bg-yellow-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-yellow-600"
                            >
                                {loading ? "Submitting…" : "Validate email"}
                            </button>
                        </div>

                        {/* resend */}
                        <div className="mt-4 text-center">
                            <button
                                onClick={handleResend}
                                disabled={loading}
                                className="text-sm underline text-yellow-600 hover:text-yellow-500 disabled:opacity-50"
                            >
                                {loading ? "Sending..." : "Send code again"}
                            </button>
                        </div>

                        {error && <p className="text-sm text-red-600">{error}</p>}
                    </form>

                    <p className="mt-10 text-center text-sm/6 text-gray-500">
                        {message}
                    </p>
                </div>
            </div>
        </>
    );
}
