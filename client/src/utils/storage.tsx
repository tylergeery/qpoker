import { Player } from "../objects/Player";

export class UserStorage {
    getID(): string {
        return localStorage.getItem("user_id");
    }

    getToken(): string {
        return localStorage.getItem("user_token");
    }

    setUser(player: Player) {
        localStorage.setItem("user_id", player.id.toString());
        localStorage.setItem("user_token", player.token);
    }

    removePlayer() {
        this.deleteToken();
        this.deleteID();
    }

    deleteID() {
        return localStorage.removeItem("user_id");
    }

    deleteToken() {
        return localStorage.removeItem("user_token");
    }
}

export const userStorage = new UserStorage();
