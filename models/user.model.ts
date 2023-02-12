import mongoose, { Schema, Document } from 'mongoose';

export interface UserType extends Document {
    email: string;
    password: string;
}

const UserSchema: Schema = new Schema({
    email: {
        type: String,
        required: true,
        unique: true,
    },
    password: {
        type: String,
        required: true,
    }
}, {
    timestamps: true
});

export default mongoose.model<UserType>('user', UserSchema);
