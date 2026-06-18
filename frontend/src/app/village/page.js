"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/services/api";
import Building from "@/components/building"

export default function VillagePage() {
    const [village, setVillage] = useState(null);
    const [loading, setLoading] = useState(true);
    const [selectedBuilding, setSelectedBuilding] = useState(null);
    const [now, setNow] = useState(Date.now());
    

    async function fetchVillage() {
        try {
        const res = await apiFetch("/village");
        const data = await res.json();

        if (!res.ok) {
            throw new Error(data.message || "Failed to fetch village");
        }

        setVillage(data);
        if (selectedBuilding) {
            const updated = data.buildings?.find(b => b.id === selectedBuilding.id);
            if (updated) setSelectedBuilding(updated);
        }
        } catch (err) {
            console.error(err);
            alert("Failed to load village");
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        fetchVillage();
    }, []);

    useEffect(() => {
        const interval = setInterval(() => {
            setNow(Date.now());
        }, 1000);

        return () => clearInterval(interval);
    }, []);

    function handleSelectBuilding(building) {
        setSelectedBuilding(building)
    }

    async function upgradeBuilding() {
        if (!selectedBuilding) return;

        try {
            const res = await apiFetch("/village/buildings/upgrade", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    building_instance_id: selectedBuilding.id,
                }),
            });

            const text = await res.text();

            if (!res.ok) {
                throw new Error(text);
            }

            await fetchVillage();

            alert("Upgrade started");
        } catch (err) {
            console.error(err);
            alert(err.message);
        }
    }


    if (loading) {
        return <div className="p-6">Loading village...</div>;
    }

    if (!village) {
        console.log(village);
        return <div className="p-6">No village found</div>;
    }

    return (
        <div className="p-6 space-y-6">

            <div className="border p-4">
                <h2>Resources</h2>
                <p>Gold: {village.gold}</p>
                <p>Elixir: {village.elixir}</p>
            </div>

            <div className="border p-4">
                <h2>Townhall</h2>
                <p>Level: {village.townhall_level}</p>
            </div>

            <div className="border p-4">
                <h2 className="mb-4">Village Layout</h2>

                <div className="relative w-[600px] h-[600px] border bg-[#3F704D]">

                    {village.buildings?.map((b) => (
                        <Building
                            key={b.id}
                            building={b}
                            isSelected={selectedBuilding?.id===b.id}
                            onSelect={handleSelectBuilding}
                        />
                    ))}

                </div>
            </div>

            {selectedBuilding && (
                <div className="border p-4 mt-4">
                    <h2>Building Details</h2>
                    <p>ID: {selectedBuilding.id}</p>
                    <p>Level: {selectedBuilding.level}</p>

                    {(selectedBuilding.building_id === 1 && selectedBuilding.level >= 4) || 
                    (selectedBuilding.building_id !== 1 && selectedBuilding.level >= village.townhall_level) ? (
                        <button className="border px-4 py-2 mt-2 bg-gray-300 text-gray-500 cursor-not-allowed" disabled>
                            Max Level Reached
                        </button>
                    ) : (
                        <button 
                            onClick={upgradeBuilding}
                            className="border px-4 py-2 mt-2 bg-blue-500 text-white"
                        >
                            Upgrade
                        </button>
                    )}
                </div>
            )}

        </div>
    );



}