import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgStakeServicer } from "./types/poktroll/servicer/tx";
import { MsgClaim } from "./types/poktroll/servicer/tx";
import { MsgProof } from "./types/poktroll/servicer/tx";
import { MsgUnstakeServicer } from "./types/poktroll/servicer/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/poktroll.servicer.MsgStakeServicer", MsgStakeServicer],
    ["/poktroll.servicer.MsgClaim", MsgClaim],
    ["/poktroll.servicer.MsgProof", MsgProof],
    ["/poktroll.servicer.MsgUnstakeServicer", MsgUnstakeServicer],
    
];

export { msgTypes }