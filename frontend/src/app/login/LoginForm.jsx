"use client";

import { useState } from "react";
import { apiFetch } from "@/services/api";
import { useAuth } from "@/context/auth_context";
import { useRouter } from "next/navigation";
import Link from "next/link";

export default function LoginForm() {
  const { login } = useAuth();
  const router = useRouter();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  async function handleSubmit(e) {
    e.preventDefault();

    try {
      const response = await apiFetch("/login", {
        method: "POST",
        body: JSON.stringify({ username, password }),
      });

      const text = await response.text();

      if (!response.ok) {
        alert(text || "Login failed");
        return;
      }

      const data = JSON.parse(text);
      login(data.token);
      router.push("/village");
    } catch (err) {
      console.error(err);
      alert("Server error");
    }
  }

  return (
    <>
      <form
        onSubmit={handleSubmit}
        className="flex flex-col gap-4"
      >
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          className="border p-2 rounded"
        />

        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="border p-2 rounded"
        />

        <button
          type="submit"
          className="border p-2 rounded bg-black text-white"
        >
          Login
        </button>
      </form>
      <p className="mt-4">
        New User?{" "}
        <Link href="/register" className="underline">
          Register
        </Link>
      </p>
    </>

    
  );
}