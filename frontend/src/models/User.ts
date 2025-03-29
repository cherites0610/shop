class User {
    userId: string;
    displayName: string;
    picutureUrl: string;

    constructor(userId: string, displayName: string, picutureUrl: string) {
        this.userId = userId;
        this.displayName = displayName
        this.picutureUrl = picutureUrl
    }
}

export default User;
