// Code generated by fastssz. DO NOT EDIT.
// Hash: bbaeb0256e64db3d2b0744be161925a5f4aa9c9e9716d5f5e4f80984a9c7995f
// Version: 0.1.3
package types

import (
	ssz "github.com/ferranbt/fastssz"
)

// MarshalSSZ ssz marshals the BlobSidecars object
func (b *BlobSidecars) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(b)
}

// MarshalSSZTo ssz marshals the BlobSidecars object to a target array
func (b *BlobSidecars) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(4)

	// Offset (0) 'Sidecars'
	dst = ssz.WriteOffset(dst, offset)

	// Field (0) 'Sidecars'
	if size := len(b.Sidecars); size > 6 {
		err = ssz.ErrListTooBigFn("BlobSidecars.Sidecars", size, 6)
		return
	}
	for ii := 0; ii < len(b.Sidecars); ii++ {
		if dst, err = b.Sidecars[ii].MarshalSSZTo(dst); err != nil {
			return
		}
	}

	return
}

// UnmarshalSSZ ssz unmarshals the BlobSidecars object
func (b *BlobSidecars) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 4 {
		return ssz.ErrSize
	}

	tail := buf
	var o0 uint64

	// Offset (0) 'Sidecars'
	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
		return ssz.ErrOffset
	}

	if o0 < 4 {
		return ssz.ErrInvalidVariableOffset
	}

	// Field (0) 'Sidecars'
	{
		buf = tail[o0:]
		num, err := ssz.DivideInt2(len(buf), 131544, 6)
		if err != nil {
			return err
		}
		b.Sidecars = make([]*BlobSidecar, num)
		for ii := 0; ii < num; ii++ {
			if b.Sidecars[ii] == nil {
				b.Sidecars[ii] = new(BlobSidecar)
			}
			if err = b.Sidecars[ii].UnmarshalSSZ(buf[ii*131544 : (ii+1)*131544]); err != nil {
				return err
			}
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the BlobSidecars object
func (b *BlobSidecars) SizeSSZ() (size int) {
	size = 4

	// Field (0) 'Sidecars'
	size += len(b.Sidecars) * 131544

	return
}

// HashTreeRoot ssz hashes the BlobSidecars object
func (b *BlobSidecars) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(b)
}

// HashTreeRootWith ssz hashes the BlobSidecars object with a hasher
func (b *BlobSidecars) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'Sidecars'
	{
		subIndx := hh.Index()
		num := uint64(len(b.Sidecars))
		if num > 6 {
			err = ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range b.Sidecars {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 6)
	}

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the BlobSidecars object
func (b *BlobSidecars) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(b)
}
