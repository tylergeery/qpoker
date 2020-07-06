export const classNames = (...potentials: any[]): string => {
    let approved = [];

    for (let i=0; i < potentials.length; i++) {
        if (typeof potentials[i] == 'string') {
            approved.push(potentials[i]);
            continue;
        }

        for (let key in potentials[i]) {
            if (potentials[i][key]) {
                approved.push(key);
            }
        }
    }

    return approved.join(" ");
}

export type ClientAction = {
    type: string;
    data: any;
}

export const createAdminAction = (data: any): ClientAction => ({type: 'admin', data });
export const createGameAction = (data: any): ClientAction => ({type: 'game', data });
export const createChatAction = (chat: string): ClientAction => ({type: 'chat', data: chat });
export const createVideoAction = (data: any): ClientAction => ({type: 'video', data });
