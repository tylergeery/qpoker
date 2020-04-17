export interface RequestParams {
    id?: string;
    userToken?: string;
    data?: object;
}

export abstract class BaseRequest<T> {
    success: boolean;
    errors?: string[];

    abstract getEndpoint(params: RequestParams): string;

    public getURL(params: RequestParams): string {
        return window.QPoker.hostname + '/' + this.getEndpoint(params);
    }

    public getJSONHeaders(params: RequestParams): any {
        let headers: any = {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        };

        if (params.userToken) {
            headers['Authorization'] = `Bearer ${params.userToken}`;
        }

        return headers;
    }

    public async handleResponse(response: Response): Promise<T> {
        this.success = response.ok;
        let json = await response.json();

        if (this.success) {
            return json;
        }

        this.errors = json['errors'];

        return null;
    }

    public abstract async request(params: RequestParams): Promise<T>;
}