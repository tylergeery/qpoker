import { RequestParams, BaseRequest } from "./base";

export class GameHistoryRequest<T> extends BaseRequest<T> {
    public getEndpoint(params: RequestParams): string {
        return `api/v1/games/${params.id}/history`;
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