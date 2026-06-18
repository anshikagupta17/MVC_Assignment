"use client";

import { useState } from "react";
import { apiFetch } from "@/services/api";
import { useAuth } from "@/context/auth_context";

export default function LoginPage() {
  const { login } = useAuth();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  async function handleSubmit(e) {
    e.preventDefault();

    try {
      const response = await apiFetch("/login", {
        method: "POST",
        body: JSON.stringify({
          username,
          password,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        alert(data.message || "Login failed");
        return;
      }

      login(data.token);

      alert("Login successful");
    } catch (err) {
      console.error(err);
      alert("Server error");
    }
  }

  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-6">
        Login
      </h1>

      <form
        onSubmit={handleSubmit}
        className="flex flex-col gap-4 max-w-sm"
      >
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) =>
            setUsername(e.target.value)
          }
          className="border p-2"
        />

        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) =>
            setPassword(e.target.value)
          }
          className="border p-2"
        />

        <button
          type="submit"
          className="border p-2"
        >
          Login
        </button>
      </form>
    </div>
  );
}