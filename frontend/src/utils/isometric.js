const TILE_WIDTH = 64;
const TILE_HEIGHT = 32;

export function tileToScreen(tileX, tileY) {
    const screenX = (tileX - tileY) * (TILE_WIDTH / 2);
    const screenY = (tileX + tileY) * (TILE_HEIGHT / 2);
    return { screenX, screenY };
}

export function screenToTile(screenX, screenY) {
    const tileX = (screenX / (TILE_WIDTH / 2) + screenY / (TILE_HEIGHT / 2)) / 2;
    const tileY = (screenY / (TILE_HEIGHT / 2) - screenX / (TILE_WIDTH / 2)) / 2;
    return {
        tileX: Math.round(tileX),
        tileY: Math.round(tileY),
    };
}

export { TILE_WIDTH, TILE_HEIGHT };