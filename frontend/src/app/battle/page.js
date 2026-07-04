"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/services/api";
import Link from "next/link";

export default function BattlePage() {
    const [opponent, setOpponent] = useState(null);
    const [loading, setLoading] = useState(false);
    const [battleResult, setBattleResult] = useState(null);
    const [troops, setTroops] = useState([]);
    const [selectedTroops, setSelectedTroops] = useState({});
    const [battlePhase, setBattlePhase] = useState("idle");
    const [battleMessage, setBattleMessage] = useState("");
    const [messageIndex, setMessageIndex] = useState(0);
    const [errorMsg, setErrorMsg] = useState("");

    const battleMessages = [
        "Troops advancing",
        "Defenses fighting back",
        "Breaking through",
        "Townhall in DANGER",
        "Final assault underway",
    ];


    useEffect(() => {
        fetchTroops();
    }, []);

    async function fetchTroops() {
        try {
            const res = await apiFetch("/village/troops");
            const data = await res.json();

            if (!res.ok) {
                throw new Error(data.message || "Failed to load troops");
            }

            setTroops(data || []);
        } catch (err) {
            console.error(err);
        }
    }

    function showError(rawMessage, fallbackMessage) {
        const displayMessage = fallbackMessage || rawMessage;
        setErrorMsg(displayMessage);
        setTimeout(() => setErrorMsg(""), 5000);
    }


    async function FindOpponent() {
        setLoading(true);
        setOpponent(null);
        setBattleResult(null);

        try {
            const res = await apiFetch("/battle/matchmake", {
                method: "POST",
            });

            const data = await res.json();

            if (!res.ok) {
                throw new Error(data.message || "No opponent found");
            }

            setOpponent(data);
        } catch (err) {
            console.error(err);
            showError(err.message, "Couldn't find an opponent right now. Try again in a bit.");
        } finally {
            setLoading(false);
        }
    }

    function handleTroopQuantityChange(troopId, quantity) {
        setSelectedTroops((prev) => ({
            ...prev,
            [troopId]: quantity,
        }));
    }

    function buildTroopsToSend(selectedTroops) {
        const troopsToSend = [];

        for (const troopIdStr in selectedTroops) {
            const quantity = selectedTroops[troopIdStr];

            if (quantity > 0) {
                troopsToSend.push({
                    troop_id: Number(troopIdStr),
                    quantity: quantity,
                });
            }
        }

        return troopsToSend;
    }

    function advanceBattleMessage() {
        setMessageIndex(function nextIndex(prevIndex) {
            const newIndex = (prevIndex + 1) % battleMessages.length;
            setBattleMessage(battleMessages[newIndex]);
            return newIndex;
        });
    }

    function waitMilliseconds(ms) {
        return new Promise(function startTimer(resolve) {
            setTimeout(resolve, ms);
        });
    }

    async function attackOpponent() {
        if (!opponent) return;

        const troopsToSend = buildTroopsToSend(selectedTroops);

        if (troopsToSend.length === 0) {
            showError(null, "Pick at least one troop before attacking.");
            return;
        }

        setBattlePhase("fighting");
        setMessageIndex(0);
        setBattleMessage(battleMessages[0]);

        const messageInterval = setInterval(advanceBattleMessage, 6000);

        try {
            const res = await apiFetch("/battle/attack", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    defender_id: opponent.village_id,
                    troops: troopsToSend,
                }),
            });

            const data = await res.json();

            if (!res.ok) {
                throw new Error(data.message || "Attack failed");
            }

            await waitMilliseconds(30000);

            clearInterval(messageInterval);
            setBattleResult(data);
            setBattlePhase("result");
            setOpponent(null);
            setSelectedTroops({});
            await fetchTroops();
        } catch (err) {
            clearInterval(messageInterval);
            setBattlePhase("idle");
            console.error(err);
            showError(err.message, "The attack couldn't go through. Try again.");
        }
    }

    return (
        <div className="relative w-screen h-screen overflow-hidden bg-black">
            {errorMsg && (
                <div className="fixed top-6 left-1/2 -translate-x-1/2 bg-red-600 text-white px-4 py-2 rounded shadow-lg z-50">
                    {errorMsg}
                </div>
            )}
            {battlePhase === "idle" && !opponent && (
                <div className="flex flex-col items-center justify-center w-full h-full gap-6"
                    style={{ backgroundImage: `url('/assets/match.jpeg')`, 
                        backgroundSize: "cover"}}>
                    <Link href="/village">
                        <button className="px-4 py-2 rounded bg-[#3d1f00] border border-[#8B5E3C] text-[#f5c96e]">
                            Back to Village
                        </button>
                    </Link>
                    <button
                        onClick={FindOpponent}
                        disabled={loading}
                        className="bg-red-700 hover:bg-red-600 text-white text-xl font-bold px-10 py-4 rounded-xl border-2 border-red-400 shadow-lg"
                    >
                        {loading ? "Searching..." : "Find Opponent"}
                    </button>
                </div>
            )}
            {battlePhase==="idle" && opponent && (
                <>
                    <div
                        className="relative w-full h-full overflow-hidden flex items-center justify-center"
                        style={{
                            backgroundImage: `
                                linear-gradient(rgba(255,255,255,0.18) 1px, transparent 1px),
                                linear-gradient(90deg, rgba(255,255,255,0.18) 1px, transparent 1px),
                                url('/assets/grass.jpeg')
                            `,
                            backgroundSize: "20px 20px, 20px 20px",
                            backgroundColor: "#3F704D",
                            backgroundRepeat: "repeat",
                        }}
                    >
                        <div className="relative w-full h-full">
                            {opponent.buildings?.map((b) => (
                                <div
                                    key={b.id}
                                    className="absolute"
                                    style={{
                                        left: `${b.x * 20}px`,
                                        top: `${b.y * 20}px`,
                                        width: `${b.size_x * 20}px`,
                                        height: `${b.size_y * 20}px`,
                                    }}
                                >
                                    <img
                                        src={`/assets/buildings/${b.building_id}.png`}
                                        alt=""
                                        style={{
                                            width: "100%",
                                            height: "100%",
                                            objectFit: "contain",
                                        }}
                                    />
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="fixed top-14 left-1/2 -translate-x-1/2 z-20">
                        <div
                            className="flex gap-4 items-center px-6 py-2 rounded bg-[#3d1f00] border border-[#8B5E3C] text-[#f5c96e]"
                        >
                            <span>TH {opponent.townhall_level}</span>
                            <span>Trophies: {opponent.trophies}</span>
                            <span>Gold: {opponent.gold}</span>
                            <span>Elixir: {opponent.elixir}</span>
                        </div>
                    </div>

                    <Link href="/village" className="fixed top-4 left-4 z-20">
                        <button className="bg-[#3d1f00] border border-[#8B5E3C] text-[#f5c96e] px-4 py-2 rounded">
                            Back
                        </button>
                    </Link>

                    <button
                        onClick={FindOpponent}
                        disabled={loading}
                        className="fixed top-4 right-4 z-20 font-bold px-4 py-2 rounded bg-[#3d1f00] border border-[#8B5E3C] text-[#f5c96e]"
                    >
                        {loading ? "..." : "New Opponent"}
                    </button>

                    <button
                        onClick={() => document.getElementById("troop-modal").showModal()}
                        className="fixed bottom-6 right-6 z-20 font-bold px-8 py-4 rounded bg-[#3d1f00] border border-[#8B5E3C] text-[#f5c96e]"
                    >
                        Select Troops & Attack
                    </button>

                    <dialog
                        id="troop-modal"
                        className="rounded-xl p-0 backdrop:bg-black/70"
                        style={{ backgroundColor: "#2a1200", border: "2px solid #8B5E3C", minWidth: "400px" }}
                    >
                        <div className="p-6">
                            <div className="flex justify-between items-center mb-4">
                                <h2 className="text-xl font-bold" style={{ color: "#f5c96e" }}>Select Troops</h2>
                                <button
                                    onClick={() => document.getElementById("troop-modal").close()}
                                    style={{
                                        backgroundColor: "#3d1f00",
                                        border: "2px solid #8B5E3C",
                                        color: "#f5c96e",
                                        padding: "4px 12px",
                                        borderRadius: "6px",
                                        cursor: "pointer",
                                    }}
                                >
                                    Close
                                </button>
                            </div>

                            <div className="grid grid-cols-2 gap-3 mb-4">
                                {troops.filter(function hasTroops(t) {
                                    return t.quantity > 0;
                                }).map((t) => (
                                    <div
                                        key={t.troop_id}
                                        className={`bg-[#3d1f00] border-2 rounded-xl p-3 flex flex-col items-center gap-1 ${
                                            selectedTroops[t.troop_id] > 0
                                                ? "border-yellow-400 shadow-yellow-400 shadow-md"
                                                : "border-[#8B5E3C]"
                                        }`}
                                    >
                                        <p className="text-[#f5c96e] font-bold text-sm tracking-wide uppercase">
                                            {t.name}
                                        </p>
                                        <img
                                            src={`/assets/troops/${t.troop_id}.png`}
                                            alt={t.name}
                                            className="w-16 h-16 object-contain"
                                        />
                                        <div className="text-xs text-yellow-300">
                                            <span>Owned: {t.quantity}</span>
                                        </div>
                                        <input
                                            type="number"
                                            min="0"
                                            max={t.quantity}
                                            value={selectedTroops[t.troop_id] || 0}
                                            onChange={(e) =>
                                                handleTroopQuantityChange(t.troop_id, Number(e.target.value))
                                            }
                                            className="w-16 text-center bg-[#1a0a00] border border-[#8B5E3C] text-white rounded px-1 py-0.5 text-sm"
                                        />
                                    </div>
                                ))}
                            </div>

                            <button
                                onClick={() => {
                                    document.getElementById("troop-modal").close();
                                    attackOpponent();
                                }}
                                className="w-full bg-green-700 hover:bg-green-600 text-white font-bold py-3 rounded-xl border-2 border-green-400 text-lg"
                            >
                                Attack
                            </button>
                        </div>
                    </dialog>
                </>
            )}

            {battlePhase === "fighting" && (
                <div className="fixed inset-0 bg-black flex flex-col items-center justify-center z-50">
                    <p className="text-4xl font-bold text-[#f5c96e] animate-pulse tracking-widest" style={{ fontFamily: "serif" }}>
                        {battleMessage}
                    </p>
                </div>
            )}

            {battlePhase === "result" && battleResult && (
                <div className="fixed inset-0 bg-black/80 flex items-center justify-center z-50">
                    <div className="max-w-md w-full bg-[#2a1200] border-2 border-[#8B5E3C] rounded-xl p-6 text-center">
                        <h2 className="text-2xl font-bold text-[#f5c96e] mb-4" style={{ fontFamily: "serif" }}>
                            Battle Result
                        </h2>

                        <div className="flex justify-center gap-2 mb-4">
                            {[1, 2, 3].map((star) => (
                                <span
                                    key={star}
                                    className={`text-4xl ${star <= battleResult.stars ? "text-yellow-400" : "text-gray-600"}`}
                                >
                                    ★
                                </span>
                            ))}
                        </div>

                        <div className="space-y-2 text-lg text-white">
                            <p>Destruction: <span className="text-[#f5c96e] font-bold">{battleResult.destruction}%</span></p>
                            <p>Gold looted: <span className="text-yellow-300 font-bold">{battleResult.loot_gold}</span></p>
                            <p>Elixir looted: <span className="text-purple-300 font-bold">{battleResult.loot_elixir}</span></p>
                            <p>Trophy change: <span className={`font-bold ${battleResult.trophy_change >= 0 ? "text-green-400" : "text-red-400"}`}>
                                {battleResult.trophy_change >= 0 ? "+" : ""}{battleResult.trophy_change}
                            </span></p>
                        </div>

                        <button
                            onClick={FindOpponent}
                            className="mt-6 bg-red-700 hover:bg-red-600 text-white font-bold px-8 py-3 rounded-xl border-2 border-red-400"
                        >
                            Battle Again
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
}