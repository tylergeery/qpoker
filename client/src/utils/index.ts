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