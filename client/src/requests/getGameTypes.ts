import { RequestParams, BaseRequest } from "./base";

export class GameTypesRequest<T> extends BaseRequest<T> {
    public getEndpoint(params: RequestParams): string {
        return `api/v1/games/types`;
    }

    public async request(params: RequestParams): Promise<T> {
        let response = await fetch(this.getURL(params), {
            method: 'GET',
            credentials: 'include',
            headers: this.getJSONHeaders(params)
        });

        return this.handleResponse(response);
    }
}