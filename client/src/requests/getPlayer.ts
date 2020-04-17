import { RequestParams, BaseRequest } from "./base";

export class GetPlayerRequest<T> extends BaseRequest<T> {
    public getEndpoint(params: RequestParams): string {
        return 'api/v1/players/' + params.id;
    }

    public async request(params: RequestParams): Promise<T> {
        let response = await fetch(this.getURL(params), {
            method: 'GET',
            credentials: 'include',
            headers: this.getJSONHeaders(params)
        });

        this.success = response.ok;

        return await response.json();
    }
}