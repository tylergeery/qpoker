import { RequestParams, BaseRequest } from "./base";

export class UpdateGameRequest<T> extends BaseRequest<T> {
    public getEndpoint(params: RequestParams): string {
        return 'api/v1/games/' + params.id;
    }

    public async request(params: RequestParams): Promise<T> {
        let response = await fetch(this.getURL(params), {
            method: 'PUT',
            body: JSON.stringify(params.data),
            credentials: 'include',
            headers: this.getJSONHeaders(params)
        });

        return this.handleResponse(response);
    }
}