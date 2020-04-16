import { Player } from "../objects/Player";

export class GetPlayerRequest {
    success: boolean

    public getURL(id: string): string {
        return 'http://localhost:8080/api/v1/players/' + id;
    }

    public getHeaders(userToken: string): any {
        return {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${userToken}`
        };
    }

    public async request(id: string, userToken: string): Promise<Player> {
        let response = await fetch(this.getURL(id), {
            method: 'GET',
            credentials: 'include',
            headers: this.getHeaders(userToken)
        });

        this.success = response.ok;

        return await response.json();
    }
}