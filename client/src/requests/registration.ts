import { Player } from "../objects/Player";

export class RegistrationRequest {
    success: boolean

    public getURL(): string {
        return 'http://localhost:8080/api/v1/players';
    }

    public getHeaders(): any {
        return {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        };
    }

    public async request(registrationData: object): Promise<Player> {
        let response = await fetch(this.getURL(), {
            method: 'POST',
            body: JSON.stringify(registrationData),
            credentials: 'include',
            headers: this.getHeaders()
        });

        this.success = response.ok;

        return await response.json();
    }
}