package homework

import (
	"sort"
	"strconv"
)

// Single Number (只出现一次的数字)
func SingleNumber(nums []int) int {
	// 使用map记录每个数字出现的次数
	m := make(map[int]int)
	for _, num := range nums {
		m[num]++
	}
	for num, count := range m {
		if count == 1 {
			return num
		}
	}
	return 0
}

// Is Palindrome (回文数)
func IsPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	// 将数字转换为字符串
	s := strconv.Itoa(x)
	for i := 0; i < len(s)/2; i++ {
		// 比较字符串的第i个字符和第len(s)-i-1个字符
		if s[i] != s[len(s)-i-1] {
			return false
		}
	}
	return true
}

// Is Valid Parentheses (有效的括号)
func IsValidParentheses(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	// 定义映射
	parentMap := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}
	// 使用栈来存储左括号
	stack := make([]rune, 0)
	// 遍历字符串
	for _, char := range s {
		if char == '(' || char == '[' || char == '{' {
			stack = append(stack, char)
		} else {
			if len(stack) == 0 {
				return false
			} else {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top != parentMap[char] {
					return false
				}
			}
		}
	}
	return len(stack) == 0
}

// Longest Common Prefix (最长公共前缀)
func LongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	// 遍历字符串数组
	for i := 1; i < len(strs); i++ {
		// 比较当前字符串和prefix的公共前缀
		for j := 0; j < len(prefix); j++ {
			if j >= len(strs[i]) || strs[i][j] != prefix[j] {
				prefix = prefix[:j]
				break
			}
		}
	}
	return prefix
}

// Plus One (加一)
func PlusOne(digits []int) []int {
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}
	digits = append([]int{1}, digits...)
	return digits
}

// Remove Duplicates (删除有序数组中的重复项)
func RemoveDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	prev := nums[0]
	count := 1
	for i := 1; i < len(nums); i++ {
		if nums[i] != prev {
			nums[count] = nums[i]
			prev = nums[i]
			count++
		}
	}
	return count
}

// Merge Intervals (合并区间)
func MergeIntervals(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return [][]int{}
	}
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	merged := [][]int{intervals[0]}
	for i := 1; i < len(intervals); i++ {
		last := merged[len(merged)-1]
		current := intervals[i]
		if current[0] <= last[1] {
			last[1] = maxInt(last[1], current[1])
		} else {
			merged = append(merged, current)
		}
	}
	return merged
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Two Sum (两数之和)
func TwoSum(nums []int, target int) []int {
	if len(nums) < 2 {
		return []int{}
	}
	numMap := make(map[int]int)
	for i, num := range nums {
		complement := target - num
		if index, ok := numMap[complement]; ok {
			return []int{index, i}
		}
		numMap[num] = i
	}
	return []int{}
}
