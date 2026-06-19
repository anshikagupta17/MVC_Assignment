"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/services/api";
import Building from "@/components/building"

export default function VillagePage() {
    const [village, setVillage] = useState(null);
    const [loading, setLoading] = useState(true);
    const [selectedBuilding, setSelectedBuilding] = useState(null);
    const [showShop, setShowShop] = useState(false);
    const [shopBuildings, setShopBuildings] = useState([]);
    const [buildingToPlace, setBuildingToPlace] = useState(null);
    const [showTroops, setShowTroops] = useState(false);
    const [troops, setTroops] = useState([]);
    const [selectedTroop, setSelectedTroop] = useState(null);
    const [trainQuantity, setTrainQuantity] = useState(1);
    

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

    function handleSelectBuilding(building) {
       if (selectedBuilding?.id===building.id) {
        setSelectedBuilding(null);
       }else {
            setSelectedBuilding(building);
        }
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
    function handleGridClick(e) {
        const rect = e.currentTarget.getBoundingClientRect();
        let x = Math.floor((e.clientX - rect.left) / 20);
        let y = Math.floor((e.clientY - rect.top) / 20);

        if (buildingToPlace) {
            const sizeX = buildingToPlace.size_x || 1;
            const sizeY = buildingToPlace.size_y || 1;
            x = Math.max(0, Math.min(x, 49 - (sizeX - 1)));
            y = Math.max(0, Math.min(y, 49 - (sizeY - 1)));
            placeBuilding(x, y);
            return;
        }

        if (!selectedBuilding) return;

        const sizeX = selectedBuilding.size_x || 1;
        const sizeY = selectedBuilding.size_y || 1;
        x = Math.max(0, Math.min(x, 49 - (sizeX - 1)));
        y = Math.max(0, Math.min(y, 49 - (sizeY - 1)));
        moveBuilding(x, y);
    }

    async function moveBuilding(x, y) {
        if (!selectedBuilding) return;

        try {
            const res = await apiFetch("/village/buildings/move", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    building_instance_id: selectedBuilding.id,
                    x,
                    y,
                }),
            });

            const text = await res.text();

            if (!res.ok) {
                throw new Error(text);
            }
            setSelectedBuilding(null)

            await fetchVillage();
        } catch (err) {
            console.error(err);
            alert(err.message);
        }
    }

    async function openShop() {
        try {
            const res = await apiFetch("/village/shop");
            const data = await res.json();

            if (!res.ok) {
                throw new Error(data.message || "Failed to load shop");
            }

            setShopBuildings(data);
            setShowShop(true);
        } catch (err) {
            console.error(err);
            alert(err.message);
        }
    }

    function handlePickBuilding(building) {
        setBuildingToPlace(building);
        setShowShop(false);
    }

    async function placeBuilding(x, y) {
        try {
            const res = await apiFetch("/village/buildings/build", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    building_id: buildingToPlace.building_id,
                    x,
                    y,
                }),
            });

            const text = await res.text();

            if (!res.ok) {
                throw new Error(text);
            }

            setBuildingToPlace(null);
            await fetchVillage();
        } catch (err) {
            console.error(err);
            alert(err.message);
        }
    }

    async function collectResources() {
        try {
            const res = await apiFetch("/village/collect", {
                method: "POST",
            });

            const data = await res.json();

            if (!res.ok) {
                throw new Error(data.message || "Failed to collect resources");
            }

            await fetchVillage();
            alert(`Collected ${data.Gold} gold and ${data.Elixir} elixir`);
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
                <button
                    onClick={openShop}
                    className="border px-4 py-2"
                >
                    Shop
                </button>

                {buildingToPlace && (
                    <p className="mt-2 text-sm">
                        Click on the map to place: {buildingToPlace.name}
                    </p>
                )}
            </div>

            <div className="border p-4">
                <h2 className="mb-4">Village Layout</h2>

                <div className="relative w-[600px] h-[600px] border bg-[#3F704D]"
                onClick={handleGridClick}
                >

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

                    {(selectedBuilding.building_id === 6 || selectedBuilding.building_id === 7) && (
                        <button
                            onClick={collectResources}
                            className="border px-4 py-2 mt-2 bg-yellow-500 text-white"
                        >
                            Collect
                        </button>
                    )}

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
            {showShop && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                    <div className="bg-white text-black p-6 rounded max-w-md w-full max-h-[80vh] overflow-y-auto">
                        <h2 className="text-xl font-bold mb-4">Shop</h2>

                        {shopBuildings.length === 0 ? (
                            <p>No buildings available</p>
                        ) : (
                            shopBuildings.map((b) => (
                                <div
                                    key={b.building_id}
                                    className="border p-3 mb-2 flex justify-between items-center"
                                >
                                    <div>
                                        <p className="font-semibold">{b.name}</p>
                                        <p className="text-sm">
                                            Cost: {b.cost} {b.cost_type}
                                        </p>
                                        <p className="text-sm">
                                            Owned: {b.current_count}/{b.max_quantity}
                                        </p>
                                    </div>

                                    <button
                                        onClick={() => handlePickBuilding(b)}
                                        className="border px-3 py-1 bg-blue-500 text-white"
                                    >
                                        Build
                                    </button>
                                </div>
                            ))
                        )}

                        <button
                            onClick={() => setShowShop(false)}
                            className="border px-4 py-2 mt-4"
                        >
                            Close
                        </button>
                    </div>
                </div>
            )}
        </div>
    );



}