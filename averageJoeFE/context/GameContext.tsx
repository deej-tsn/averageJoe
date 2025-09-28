import { RefObject } from "react";

type GameContextType = {
    gameWebsocketRef : RefObject<WebSocket|null>
    gameID
}