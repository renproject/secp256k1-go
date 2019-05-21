/**********************************************************************
 * Copyright (c) 2013, 2014 Pieter Wuille                             *
 * Distributed under the MIT software license, see the accompanying   *
 * file COPYING or http://www.opensource.org/licenses/mit-license.php.*
 **********************************************************************/

#ifndef SECP256K1_FIELD_REPR_H
#define SECP256K1_FIELD_REPR_H

#include <stdint.h>

typedef struct secp256k1_fe {
    /* X = sum(i=0..4, n[i]*2^(i*52)) mod p
     * where p = 2^256 - 0x1000003D1
     */
    uint64_t n[5];
} secp256k1_fe;

#endif /* SECP256K1_FIELD_REPR_H */
