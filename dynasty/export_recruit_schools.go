package dynasty

import "sort"

// arrayStoreMemberRefs returns record references stored in an array-store row.
func (f *File) arrayStoreMemberRefs(storeRef *RecordReference) []*RecordReference {
	if storeRef == nil || storeRef.TableID == 0 {
		return nil
	}
	arrTable, ok := f.GetTableByID(storeRef.TableID)
	if !ok || arrTable == nil {
		return nil
	}
	if err := arrTable.ReadRecords(); err != nil {
		return nil
	}
	rowIdx := int(storeRef.RowNumber)
	if rowIdx < 0 || rowIdx >= len(arrTable.Records) {
		return nil
	}

	var refs []*RecordReference
	for _, value := range arrTable.Records[rowIdx].Fields {
		if isEmptyRecordReference(value.Reference) {
			continue
		}
		refs = append(refs, value.Reference)
	}
	return refs
}

func recruitSchoolInterests(f *File, recruit Record, schoolTable *Table, teams teamIndexMaps) []RecruitingSchoolInterestExport {
	if schoolTable == nil {
		return nil
	}
	value, ok := recruit.Get("TopSchoolsList")
	if !ok || value.Reference == nil {
		return nil
	}

	schoolID := schoolTable.Header.TableID
	var exports []RecruitingSchoolInterestExport
	seen := make(map[int]struct{})
	for _, ref := range f.arrayStoreMemberRefs(value.Reference) {
		if ref.TableID != 0 && ref.TableID != schoolID {
			continue
		}
		row, ok := f.RecordByReference("ProspectTargetSchool", ref)
		if !ok {
			continue
		}
		interest := buildSchoolInterestExport(row, teams)
		if interest == nil {
			continue
		}
		if _, dup := seen[interest.TeamID]; dup {
			continue
		}
		seen[interest.TeamID] = struct{}{}
		exports = append(exports, *interest)
	}
	if len(exports) == 0 {
		return nil
	}
	sort.Slice(exports, func(i, j int) bool {
		if exports[i].Influence != exports[j].Influence {
			return exports[i].Influence > exports[j].Influence
		}
		return exports[i].TeamID < exports[j].TeamID
	})
	return exports
}

func topSchoolFromInterests(interests []RecruitingSchoolInterestExport) *RecruitingSchoolInterestExport {
	if len(interests) == 0 {
		return nil
	}
	top := interests[0]
	return &top
}
