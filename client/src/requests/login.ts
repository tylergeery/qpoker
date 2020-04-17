import { RequestParams, BaseRequest } from "./base";

export class LoginRequest<T> extends BaseRequest<T> {
    public getEndpoint(params: RequestParams): string {
        return 'api/v1/players/login';
    }

    public async request(params: RequestParams): Promise<T> {
        let response = await fetch(this.getURL(params), {
            method: 'POST',
            body: JSON.stringify(params.data),
            credentials: 'include',
            headers: this.getJSONHeaders(params)
        });

        return this.handleResponse(response);
    }
}